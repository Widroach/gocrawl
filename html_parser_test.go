package main

import (
	"net/url"
	"strings"
	"testing"
)

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

func Test_getFirstParagraphFromHTML(t *testing.T) {
	tests := []struct {
		name string
		html string
		want string
	}{
		// --- Main tag with paragraph inside ---
		{
			name: "main tag with simple paragraph",
			html: "<main><p>Main content</p></main>",
			want: "Main content",
		},
		{
			name: "main tag with multiple paragraphs - first wins",
			html: "<main><p>First</p><p>Second</p></main>",
			want: "First",
		},
		{
			name: "main tag with paragraph and other elements",
			html: "<main><h1>Title</h1><p>Main paragraph</p></main>",
			want: "Main paragraph",
		},
		{
			name: "main tag with nested paragraph",
			html: "<main><div><p>Nested in div</p></div></main>",
			want: "Nested in div",
		},
		{
			name: "main tag with attributes",
			html: `<main class="content"><p>Attributed main</p></main>`,
			want: "Attributed main",
		},
		{
			name: "main tag with empty paragraph",
			html: "<main><p></p></main>",
			want: "",
		},
		{
			name: "main tag with whitespace paragraph",
			html: "<main><p>   </p></main>",
			want: "",
		},
		{
			name: "main tag with paragraph with inner tags",
			html: "<main><p><strong>Bold</strong> text</p></main>",
			want: "Bold text",
		},

		// --- Main tag exists but no paragraph inside ---
		{
			name: "main exists but no p inside - fallback to first p outside",
			html: "<main><div>No paragraphs here</div></main><p>Fallback paragraph</p>",
			want: "Fallback paragraph",
		},
		{
			name: "main exists empty - fallback to first p",
			html: "<main></main><p>Outside paragraph</p>",
			want: "Outside paragraph",
		},
		{
			name: "main exists with headings only - fallback",
			html: "<main><h1>Title</h1><h2>Subtitle</h2></main><p>Fallback</p>",
			want: "Fallback",
		},
		{
			name: "main exists no p anywhere",
			html: "<main><div>No paragraphs at all</div></main>",
			want: "",
		},

		// --- No main tag - fallback to first p ---
		{
			name: "no main tag - simple paragraph",
			html: "<p>Hello World</p>",
			want: "Hello World",
		},
		{
			name: "no main tag - multiple paragraphs first wins",
			html: "<p>First</p><p>Second</p>",
			want: "First",
		},

		// --- Multiple main tags ---
		{
			name: "multiple main tags - first main wins",
			html: "<main><p>First main</p></main><main><p>Second main</p></main>",
			want: "First main",
		},

		// --- Main tag edge cases ---
		{
			name: "main tag not closed",
			html: "<main><p>Unclosed main</p>",
			want: "Unclosed main",
		},
		{
			name: "main tag with mixed case",
			html: "<MAIN><p>Mixed case main</p></MAIN>",
			want: "Mixed case main",
		},
		{
			name: "main tag deeply nested in document",
			html: "<!DOCTYPE html><html><body><main><p>Deep main paragraph</p></main></body></html>",
			want: "Deep main paragraph",
		},
		{
			name: "paragraph before main - main wins",
			html: "<p>Before main</p><main><p>Inside main</p></main>",
			want: "Inside main",
		},
		{
			name: "main with paragraph and style tag",
			html: "<main><style>.x{color:red}</style><p>Visible text</p></main>",
			want: "Visible text",
		},

		// --- Single paragraph (no main) ---
		{
			name: "paragraph with extra whitespace",
			html: "<p>   Hello World   </p>",
			want: "Hello World",
		},
		{
			name: "paragraph with extra whitespace",
			html: "<p>   Hello World   </p>",
			want: "Hello World",
		},
		{
			name: "paragraph with newlines and tabs",
			html: "<p>\n\tHello\n\tWorld\n</p>",
			want: "Hello World",
		},
		{
			name: "paragraph nested inside div",
			html: "<div><p>Nested paragraph</p></div>",
			want: "Nested paragraph",
		},
		{
			name: "paragraph with attributes",
			html: `<p class="intro">A paragraph with attributes</p>`,
			want: "A paragraph with attributes",
		},
		{
			name: "paragraph with inner nested tags",
			html: "<p><strong>Bold</strong> text and <em>italic</em> text</p>",
			want: "Bold text and italic text",
		},
		{
			name: "paragraph with link inside",
			html: `<p>Visit <a href="https://example.com">our site</a> today</p>`,
			want: "Visit our site today",
		},
		{
			name: "paragraph with image inside",
			html: `<p><img src="pic.jpg" alt="photo"> Some text</p>`,
			want: "Some text",
		},
		{
			name: "paragraph with special characters",
			html: "<p>Tom &amp; Jerry &lt;are&gt; friends</p>",
			want: "Tom & Jerry <are> friends",
		},
		{
			name: "paragraph with unicode",
			html: "<p>Héllo Wörld café</p>",
			want: "Héllo Wörld café",
		},
		{
			name: "paragraph with line break tag",
			html: "<p>Hello<br>World</p>",
			want: "HelloWorld",
		},
		{
			name: "paragraph with multiple break tags",
			html: "<p>A<br>B<br>C</p>",
			want: "ABC",
		},
		{
			name: "paragraph empty text",
			html: "<p></p>",
			want: "",
		},
		{
			name: "paragraph with only whitespace",
			html: "<p>   </p>",
			want: "",
		},

		// --- Multiple paragraphs - first wins ---
		{
			name: "multiple paragraphs - first wins",
			html: "<p>First</p><p>Second</p><p>Third</p>",
			want: "First",
		},
		{
			name: "multiple paragraphs with content in between",
			html: "<p>First</p><div>stuff</div><p>Second</p>",
			want: "First",
		},
		{
			name: "first paragraph empty, second has text",
			html: "<p></p><p>Second</p>",
			want: "",
		},

		// --- Paragraph after other elements ---
		{
			name: "paragraph after heading",
			html: "<h1>Title</h1><p>Body text</p>",
			want: "Body text",
		},
		{
			name: "paragraph after div",
			html: "<div>stuff</div><p>Paragraph here</p>",
			want: "Paragraph here",
		},
		{
			name: "paragraph deep in document",
			html: "<!DOCTYPE html><html><body><p>Deep paragraph</p></body></html>",
			want: "Deep paragraph",
		},

		// --- No paragraphs ---
		{
			name: "no p tags at all",
			html: "<h1>Only a heading</h1>",
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
			name: "only h1 h2 h3 tags",
			html: "<h1>H1</h1><h2>H2</h2><h3>H3</h3>",
			want: "",
		},
		{
			name: "only div tags",
			html: "<div>Content</div>",
			want: "",
		},

		// --- Edge cases ---
		{
			name: "paragraph with mixed case tag",
			html: "<P>Mixed Case Paragraph</P>",
			want: "Mixed Case Paragraph",
		},
		{
			name: "paragraph not closed",
			html: "<p>Unclosed paragraph",
			want: "Unclosed paragraph",
		},
		{
			name: "very long paragraph",
			html: "<p>" + strings.Repeat("word ", 500) + "</p>",
			want: strings.TrimRight(strings.Repeat("word ", 500), " "),
		},
		{
			name: "paragraph with unicode escapes",
			html: "<p>Price: 100&cent;</p>",
			want: "Price: 100¢",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstParagraphFromHTML(tt.html)
			if got != tt.want {
				t.Errorf("getFirstParagraphFromHTML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getURLsFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		htmlBody string
		baseURL  *url.URL
		want     []string
		wantErr  bool
	}{
		{
			name:     "single absolute link",
			htmlBody: `<a href="https://boot.dev/blog">Blog</a>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/blog"},
		},
		{
			name:     "single relative link",
			htmlBody: `<a href="/blog">Blog</a>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/blog"},
		},
		{
			name:     "single relative link 2",
			htmlBody: `<a href="./blog">Blog</a>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/blog"},
		},
		{
			name:     "single relative link 3",
			htmlBody: `<a href="blog">Blog</a>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/blog"},
		},
		{
			name:     "relative link resolved against base",
			htmlBody: `<a href="/blog">Blog</a>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/blog"},
		},
		{
			name: "multiple links",
			htmlBody: `
				<a href="/about">About</a>
				<a href="/blog">Blog</a>
				<a href="/contact">Contact</a>
			`,
			baseURL: parseTestURL("https://boot.dev"),
			want:    []string{"https://boot.dev/about", "https://boot.dev/blog", "https://boot.dev/contact"},
		},
		{
			name:     "link inside main content",
			htmlBody: `<main><a href="/page">Click</a></main>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/page"},
		},
		{
			name:     "no links in HTML",
			htmlBody: `<p>No links here</p>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "empty HTML",
			htmlBody: ``,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "link without href",
			htmlBody: `<a>no href</a>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "nested link in heading",
			htmlBody: `<h1><a href="/title">Title</a></h1>`,
			baseURL:  parseTestURL("https://boot.dev"),
			want:     []string{"https://boot.dev/title"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getURLsFromHTML(tt.htmlBody, tt.baseURL)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getURLsFromHTML() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getURLsFromHTML() succeeded unexpectedly")
			}
			if len(got) != len(tt.want) {
				t.Errorf("getURLsFromHTML() returned %d URLs, want %d\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getURLsFromHTML() URL[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func parseTestURL(rawURL string) *url.URL {
	u, _ := url.Parse(rawURL)
	return u
}

func Test_getImagesFromHTML(t *testing.T) {
	tests := []struct {
		name     string
		htmlBody string
		baseURL  *url.URL
		want     []string
		wantErr  bool
	}{
		{
			name:     "single absolute image",
			htmlBody: `<img src="https://example.com/photo.jpg">`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/photo.jpg"},
		},
		{
			name:     "single relative image resolved against base",
			htmlBody: `<img src="/images/photo.jpg">`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/images/photo.jpg"},
		},
		{
			name:     "multiple images",
			htmlBody: `<img src="/a.png"><img src="/b.png"><img src="/c.png">`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/a.png", "https://example.com/b.png", "https://example.com/c.png"},
		},
		{
			name:     "image inside div",
			htmlBody: `<div><img src="/nested.jpg"></div>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/nested.jpg"},
		},
		{
			name:     "image inside figure",
			htmlBody: `<figure><img src="/fig.jpg" alt="a figure"></figure>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/fig.jpg"},
		},
		{
			name:     "no images in HTML",
			htmlBody: `<p>No images here</p>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "empty HTML",
			htmlBody: ``,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "image without src attribute",
			htmlBody: `<img alt="no src">`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{},
		},
		{
			name:     "mixed links and images",
			htmlBody: `<a href="/page"><img src="/thumb.jpg"></a>`,
			baseURL:  parseTestURL("https://example.com"),
			want:     []string{"https://example.com/thumb.jpg"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := getImagesFromHTML(tt.htmlBody, tt.baseURL)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getImagesFromHTML() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getImagesFromHTML() succeeded unexpectedly")
			}
			if len(got) != len(tt.want) {
				t.Errorf("getImagesFromHTML() returned %d images, want %d\ngot:  %v\nwant: %v", len(got), len(tt.want), got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getImagesFromHTML() URL[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func Test_extractPageData(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		pageURL string
		want    PageData
	}{
		{
			name: "full page with h1 and paragraph",
			html: `
				<h1>My Page Title</h1>
				<p>This is the first paragraph with some content.</p>
			`,
			pageURL: "https://example.com/page1",
			want: PageData{
				URL:            "https://example.com/page1",
				Heading:        "My Page Title",
				FirstParagraph: "This is the first paragraph with some content.",
			},
		},
		{
			name: "full page with h2 fallback",
			html: `
				<h2>Sub Heading</h2>
				<p>Some paragraph text.</p>
			`,
			pageURL: "https://example.com/page2",
			want: PageData{
				URL:            "https://example.com/page2",
				Heading:        "Sub Heading",
				FirstParagraph: "Some paragraph text.",
			},
		},
		{
			name: "page with outgoing links",
			html: `
				<h1>Links Page</h1>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
				<a href="https://external.com">External</a>
			`,
			pageURL: "https://example.com",
			want: PageData{
				URL:            "https://example.com",
				Heading:        "Links Page",
				OutgoingLinks:  []string{"https://example.com/about", "https://example.com/contact", "https://external.com"},
			},
		},
		{
			name: "page with images",
			html: `
				<h1>Gallery</h1>
				<img src="/images/photo1.jpg">
				<img src="/images/photo2.png">
			`,
			pageURL: "https://example.com/gallery",
			want: PageData{
				URL:            "https://example.com/gallery",
				Heading:        "Gallery",
				ImageURLs:      []string{"https://example.com/images/photo1.jpg", "https://example.com/images/photo2.png"},
			},
		},
		{
			name: "page with everything",
			html: `
				<h1>Full Page</h1>
				<p>First paragraph here.</p>
				<a href="/link1">Link 1</a>
				<img src="/img1.jpg">
			`,
			pageURL: "https://example.com/full",
			want: PageData{
				URL:            "https://example.com/full",
				Heading:        "Full Page",
				FirstParagraph: "First paragraph here.",
				OutgoingLinks:  []string{"https://example.com/link1"},
				ImageURLs:      []string{"https://example.com/img1.jpg"},
			},
		},
		{
			name: "empty page",
			html: ``,
			pageURL: "https://example.com/empty",
			want: PageData{
				URL: "https://example.com/empty",
			},
		},
		{
			name: "page with only paragraphs",
			html: `
				<p>First paragraph.</p>
				<p>Second paragraph.</p>
			`,
			pageURL: "https://example.com/text",
			want: PageData{
				URL:            "https://example.com/text",
				FirstParagraph: "First paragraph.",
			},
		},
		{
			name: "page with only images",
			html: `
				<img src="/a.jpg">
				<img src="/b.jpg">
			`,
			pageURL: "https://example.com/imgs",
			want: PageData{
				URL:       "https://example.com/imgs",
				ImageURLs: []string{"https://example.com/a.jpg", "https://example.com/b.jpg"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractPageData(tt.html, tt.pageURL)
			if got.URL != tt.want.URL {
				t.Errorf("URL = %v, want %v", got.URL, tt.want.URL)
			}
			if got.Heading != tt.want.Heading {
				t.Errorf("Heading = %v, want %v", got.Heading, tt.want.Heading)
			}
			if got.FirstParagraph != tt.want.FirstParagraph {
				t.Errorf("FirstParagraph = %v, want %v", got.FirstParagraph, tt.want.FirstParagraph)
			}
			if len(got.OutgoingLinks) != len(tt.want.OutgoingLinks) {
				t.Errorf("OutgoingLinks = %v (len %d), want %v (len %d)", got.OutgoingLinks, len(got.OutgoingLinks), tt.want.OutgoingLinks, len(tt.want.OutgoingLinks))
			} else {
				for i := range got.OutgoingLinks {
					if got.OutgoingLinks[i] != tt.want.OutgoingLinks[i] {
						t.Errorf("OutgoingLinks[%d] = %v, want %v", i, got.OutgoingLinks[i], tt.want.OutgoingLinks[i])
					}
				}
			}
			if len(got.ImageURLs) != len(tt.want.ImageURLs) {
				t.Errorf("ImageURLs = %v (len %d), want %v (len %d)", got.ImageURLs, len(got.ImageURLs), tt.want.ImageURLs, len(tt.want.ImageURLs))
			} else {
				for i := range got.ImageURLs {
					if got.ImageURLs[i] != tt.want.ImageURLs[i] {
						t.Errorf("ImageURLs[%d] = %v, want %v", i, got.ImageURLs[i], tt.want.ImageURLs[i])
					}
				}
			}
		})
	}
}
