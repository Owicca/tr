package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	customtemplate "html/template"

	"github.com/owicca/tr/internal/models/logs"
	"upspin.io/errors"
)

func MergeMaps(m1, m2 map[string]any) map[string]any {
	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}

func MergeMapsInterface(m1 map[string]any, m2 map[any]any) map[string]any {
	for k, v := range m2 {
		m1[k.(string)] = v
	}

	return m1
}

func MergeMapsInterfaceReverse(m1 map[any]any, m2 map[string]any) map[any]any {
	for k, v := range m2 {
		m1[k] = v
	}

	return m1
}

func b2sSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func b2sIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func timestampToUTC(timestamp int64) string {
	t := time.Unix(timestamp, 0)

	return t.Format(time.RFC3339)
}

func timestampToCustomDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)

	return fmt.Sprintf("%d-%d-%d|%d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}

func lastInArray(a []int) int {
	if len(a) == 0 {
		return -1
	}
	return a[len(a)-1]
}

func objectToJSON(obj any) string {
	const op errors.Op = "infra.template.asjson"
	var (
		results string
		buff    bytes.Buffer
	)

	js, err := json.Marshal(obj)
	if err != nil {
		logs.LogErr(op, err)
		return results
	}
	if err = json.Indent(&buff, js, "", " "); err != nil {
		logs.LogErr(op, err)
		return results
	}
	results = string(buff.String())

	return results
}

func stringToHTML(html string) customtemplate.HTML {
	return customtemplate.HTML(html)
}

func generateDict(values ...any) map[string]any {
	const op errors.Op = "utils.generateDict"

	if len(values)%2 != 0 {
		logs.LogErr(op, errors.Str("'params' should be called with pairs of values"))
		return nil
	}

	dict := make(map[string]any, len(values))
	for i := 0; i < len(values); i += 2 {
		k, ok := values[i].(string)
		if !ok {
			logs.LogErr(op, errors.Errorf("%d th key is not a string", i/2))
			return nil
		}
		dict[k] = values[i+1]
	}

	return dict
}

func GetStaticDir() string {
	wd, _ := os.Getwd()
	return fmt.Sprintf("%s/static/media", wd)
}

func GeneratePagination(total int, limit int) (int, []any) {
	if total < 2 || total < limit {
		return 1, []any{1}
	}
	pageCount := total / limit

	//if (total % limit) != 0 {
	//	pageCount += 1
	//}

	pageHelper := make([]any, pageCount)
	for idx, _ := range pageHelper {
		pageHelper[idx] = idx + 1
	}

	return pageCount, pageHelper
}

func Decrement(i int) int {
	return i - 1
}

func Increment(i int) int {
	return i + 1
}

func Contains[T comparable](haystack []T, needle T) bool {
	results := false

	for _, elem := range haystack {
		if elem == needle{
			results = true
			break;
		}
	}

	return results
}
