package main

import "testing"

func Test_getHeadingFromHTML(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		// --- H1 present (primary) ---
		{
			name: "simple h1",
			html: "<h1>Hello World</h1>",
			want: "Hello World",
		},
		{
			name: "h1 with extra whitespace",
			html: "<h1>   Hello World   </h1>",
			want: "Hello World",
		},
		{
			name: "h1 with inner nested tags",
			html: "<h1><span>Hello</span> World</h1>",
			want: "Hello World",
		},
		{
			name: "h1 nested inside div",
			html: "<div><h1>Title Here</h1></div>",
			want: "Title Here",
		},
		{
			name: "h1 with attributes",
			html: `<h1 class="main-title">Title Here</h1>`,
			want: "Title Here",
		},
		{
			name: "h1 with newlines and tabs",
			html: "<h1>\n\tHello\n\tWorld\n</h1>",
			want: "Hello World",
		},
		{
			name: "h1 empty text",
			html: "<h1></h1>",
			want: "",
		},
		{
			name: "h1 with only whitespace",
			html: "<h1>   </h1>",
			want: "",
		},
		{
			name: "h1 with special characters",
			html: "<h1>Hello &amp; World</h1>",
			want: "Hello & World",
		},
		{
			name: "h1 with unicode",
			html: "<h1>Hello \u00e9\u00e8\u00ea</h1>",
			want: "Hello \u00e9\u00e8\u00ea",
		},
		{
			name: "h1 with link inside",
			html: `<h1><a href="/page">Click Here</a></h1>`,
			want: "Click Here",
		},
		{
			name: "h1 is first element in document",
			html: "<!DOCTYPE html><html><body><h1>First</h1><h1>Second</h1></body></html>",
			want: "First",
		},
		{
			name: "h1 after other elements",
			html: "<p>Paragraph</p><h1>Title</h1>",
			want: "Title",
		},
		{
			name: "multiple h1 elements - first wins",
			html: "<h1>First</h1><p>Middle</p><h1>Second</h1>",
			want: "First",
		},

		// --- No h1, fallback to h2 ---
		{
			name: "no h1, simple h2 fallback",
			html: "<h2>Sub Heading</h2>",
			want: "Sub Heading",
		},
		{
			name: "no h1, h2 with extra whitespace",
			html: "<h2>   Sub Heading   </h2>",
			want: "Sub Heading",
		},
		{
			name: "no h1, h2 nested inside div",
			html: "<div><h2>Sub Title</h2></div>",
			want: "Sub Title",
		},
		{
			name: "no h1, h2 with attributes",
			html: `<h2 class="subtitle">Sub Title</h2>`,
			want: "Sub Title",
		},
		{
			name: "no h1, h2 with inner nested tags",
			html: "<h2><em>Important</em> Subtitle</h2>",
			want: "Important Subtitle",
		},
		{
			name: "no h1, h2 with newlines",
			html: "<h2>\n\tSub\n\tTitle\n</h2>",
			want: "Sub Title",
		},
		{
			name: "no h1, multiple h2 first wins",
			html: "<h2>First</h2><h2>Second</h2>",
			want: "First",
		},
		{
			name: "no h1, h2 empty text",
			html: "<h2></h2>",
			want: "",
		},
		{
			name: "no h1, h2 with special characters",
			html: "<h2>Tom &amp; Jerry</h2>",
			want: "Tom & Jerry",
		},

		// --- Neither h1 nor h2 ---
		{
			name: "no h1 or h2",
			html: "<h3>Only H3</h3>",
			want: "",
		},
		{
			name: "only p tags",
			html: "<p>Just a paragraph</p>",
			want: "",
		},
		{
			name: "empty string",
			html: "",
			want: "",
		},
		{
			name: "plain text no tags",
			html: "Just plain text",
			want: "",
		},
		{
			name: "only h3 h4 h5 h6",
			html: "<h3>H3</h3><h4>H4</h4><h5>H5</h5><h6>H6</h6>",
			want: "",
		},

		// --- Both h1 and h2 present ---
		{
			name: "both h1 and h2 present - h1 wins",
			html: "<h1>Main Title</h1><h2>Sub Title</h2>",
			want: "Main Title",
		},
		{
			name: "h1 and h2 with content in between",
			html: "<h1>Title</h1><p>Some text</p><h2>Subtitle</h2>",
			want: "Title",
		},

		// --- Edge cases ---
		{
			name: "h1 with mixed case tag",
			html: "<H1>Mixed Case H1</H1>",
			want: "Mixed Case H1",
		},
		{
			name: "h1 not closed",
			html: "<h1>Unclosed heading",
			want: "Unclosed heading",
		},
		{
			name: "h1 with line break tag",
			html: "<h1>Hello<br>World</h1>",
			want: "HelloWorld",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getHeadingFromHTML(tt.html)
			if got != tt.want {
				t.Errorf("getHeadingFromHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}
