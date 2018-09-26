package News

import (
	"fmt"
	"testing"
)

func TestFixImageUrls(t *testing.T) {
	tables := []struct {
		body     string
		expected string
	}{
		{
			`<img src="/img/krowa.jpg" width=800>`,
			fmt.Sprintf(`<img src="%s/img/krowa.jpg" width=800>`, baseUrl),
		},
		{
			`<p>test</p><img src="/img/krowa.jpg" width=800>`,
			fmt.Sprintf(`<p>test</p><img src="%s/img/krowa.jpg" width=800>`, baseUrl),
		},
	}

	for _, table := range tables {
		result := fixImageUrls(table.body)
		if result != table.expected {
			t.Errorf(`Wrong result. Got "%s", expected: "%s"`, result, table.expected)
		}
	}

}

func TestCleanUpTitle(t *testing.T) {
	tables := []struct {
		title    string
		expected string
	}{
		{
			`Mecz żużlowy w dniu 16.09.2018.`,
			`Mecz żużlowy w dniu 16.09.2018`,
		},
		{
			`11.09.2018r. - mecz piłkarski na Stadionie Wrocław.`,
			`Mecz piłkarski na Stadionie Wrocław`,
		},
		{
			`9 i 16.09.2018r. - Festiwal Wratislavia Cantans.`,
			`Festiwal Wratislavia Cantans`,
		},
	}

	for _, table := range tables {
		result := cleanUpTitle(table.title)
		if result != table.expected {
			t.Errorf(`Wrong result. Got "%s", expected: "%s"`, result, table.expected)
		}
	}
}
