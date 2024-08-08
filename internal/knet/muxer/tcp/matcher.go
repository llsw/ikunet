package tcp

import (
	"fmt"
	"strings"
	"time"

	"github.com/llsw/ikunet/internal/knet/balance"
	cmap "github.com/orcaman/concurrent-map"
)

type cacheData struct {
	vt   time.Time
	data map[string]struct{}
}

var cache cmap.ConcurrentMap = cmap.New()
var tcpFuncs = map[string]func(*matchersTree, ...string) error{
	"Uuid":   expectParameter(uuid),
	"Verson": expectParameter(version),
}

func expectParameter(fn func(*matchersTree, ...string) error) func(*matchersTree, ...string) error {
	return func(route *matchersTree, s ...string) error {
		if len(s) != 1 {
			return fmt.Errorf("unexpected number of parameters; got %d, expected 1", len(s))
		}
		return fn(route, s...)
	}
}

func stringMatcher(tree *matchersTree, s []string, prefix string, getElm func(meta Data) string) error {
	if len(s) != 1 {
		return fmt.Errorf("string matcher %s unexpected number of parameters; got %d, expected 1", prefix, len(s))
	}
	if s[0] == "*" {
		tree.matcher = func(meta Data) bool {
			return true
		}
		return nil
	}

	str := s[0]

	key := fmt.Sprintf("%s:%s", prefix, str)
	var inMap map[string]struct{}

	if cd, ok := cache.Get(key); ok {
		tempData := cd.(*cacheData)
		inMap = tempData.data
		tempData.vt = time.Now()
	} else {
		elems := strings.Split(str, "|")
		ll := len(elems)
		inMap = make(map[string]struct{}, ll)
		for _, v := range elems {
			inMap[v] = struct{}{}
		}
		cache.Set(key, &cacheData{
			vt:   time.Now(),
			data: inMap,
		})
	}

	tree.matcher = func(meta Data) bool {
		elem := getElm(meta)
		if _, ok := inMap[elem]; ok {
			return true
		}
		return false
	}
	return nil
}

func uuid(tree *matchersTree, s ...string) error {
	return stringMatcher(tree, s, "uuid", func(meta Data) string {
		return meta.req.Meta.GetUuid()
	})
}

func version(tree *matchersTree, s ...string) error {
	return stringMatcher(tree, s, "version", func(meta Data) string {
		return balance.GetTagVal(meta.instance, balance.TAG_VERSION)
	})
}

// func clientIP(tree *matchersTree, clientIP ...string) error {
// 	checker, err := ip.NewChecker(clientIP)
// 	if err != nil {
// 		return fmt.Errorf("initializing IP checker for ClientIP matcher: %w", err)
// 	}

// 	tree.matcher = func(meta Data) bool {
// 		ok, err := checker.Contains(meta.remoteIP)
// 		if err != nil {
// 			log.Warn().Err(err).Msg("ClientIP matcher: could not match remote address")
// 			return false
// 		}
// 		return ok
// 	}

// 	return nil
// }

// var hostOrIP = regexp.MustCompile(`^[[:alnum:]\.\-\:]+$`)

// // hostSNI checks if the SNI Host of the connection match the matcher host.
// func hostSNI(tree *matchersTree, hosts ...string) error {
// 	host := hosts[0]

// 	if host == "*" {
// 		// Since a HostSNI(`*`) rule has been provided as catchAll for non-TLS TCP,
// 		// it allows matching with an empty serverName.
// 		tree.matcher = func(meta Data) bool { return true }
// 		return nil
// 	}

// 	if !hostOrIP.MatchString(host) {
// 		return fmt.Errorf("invalid value for HostSNI matcher, %q is not a valid hostname", host)
// 	}

// 	tree.matcher = func(meta Data) bool {
// 		if meta.serverName == "" {
// 			return false
// 		}

// 		if host == meta.serverName {
// 			return true
// 		}

// 		// trim trailing period in case of FQDN
// 		host = strings.TrimSuffix(host, ".")

// 		return host == meta.serverName
// 	}

// 	return nil
// }

// // hostSNIRegexp checks if the SNI Host of the connection matches the matcher host regexp.
// func hostSNIRegexp(tree *matchersTree, templates ...string) error {
// 	template := templates[0]

// 	if !isASCII(template) {
// 		return fmt.Errorf("invalid value for HostSNIRegexp matcher, %q is not a valid hostname", template)
// 	}

// 	re, err := regexp.Compile(template)
// 	if err != nil {
// 		return fmt.Errorf("compiling HostSNIRegexp matcher: %w", err)
// 	}

// 	tree.matcher = func(meta Data) bool {
// 		return re.MatchString(meta.serverName)
// 	}

// 	return nil
// }

// // isASCII checks if the given string contains only ASCII characters.
// func isASCII(s string) bool {
// 	for i := range len(s) {
// 		if s[i] >= utf8.RuneSelf {
// 			return false
// 		}
// 	}

// 	return true
// }
