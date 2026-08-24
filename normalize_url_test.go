package main

import "testing"

func Test_normalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		want    string
		wantErr bool
	}{
		{
			name:   "remove https scheme",
			rawURL: "https://www.boot.dev/blog/path",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "remove http scheme",
			rawURL: "http://www.boot.dev/blog/path",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:    "invalid url",
			rawURL:  "asd://:www.boot.dev/blog/path",
			want:    "",
			wantErr: true,
		},
		{
			name:    "without schema",
			rawURL:  "www.boot.dev/blog/path",
			want:    "www.boot.dev/blog/path",
			wantErr: false,
		},
		{
			name:   "https with port",
			rawURL: "https://www.boot.dev:8080/blog/path",
			want:   "www.boot.dev:8080/blog/path",
		},
		{
			name:   "http with default port 80",
			rawURL: "http://www.boot.dev:80/blog/path",
			want:   "www.boot.dev:80/blog/path",
		},
		{
			name:   "https with default port 443",
			rawURL: "https://www.boot.dev:443/blog/path",
			want:   "www.boot.dev:443/blog/path",
		},
		{
			name:   "with query string",
			rawURL: "https://www.boot.dev/blog/path?foo=bar&baz=qux",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "with fragment",
			rawURL: "https://www.boot.dev/blog/path#section1",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "with query and fragment",
			rawURL: "https://www.boot.dev/blog/path?foo=bar#section1",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "root path only",
			rawURL: "https://www.boot.dev/",
			want:   "www.boot.dev/",
		},
		{
			name:   "no path at all",
			rawURL: "https://www.boot.dev",
			want:   "www.boot.dev",
		},
		{
			name:   "deeply nested path",
			rawURL: "https://www.boot.dev/a/b/c/d/e/f/g",
			want:   "www.boot.dev/a/b/c/d/e/f/g",
		},
		{
			name:   "path with trailing slash",
			rawURL: "https://www.boot.dev/blog/path/",
			want:   "www.boot.dev/blog/path/",
		},
		{
			name:   "path with multiple consecutive slashes",
			rawURL: "https://www.boot.dev/blog//path///sub",
			want:   "www.boot.dev/blog//path///sub",
		},
		{
			name:   "with userinfo",
			rawURL: "https://user:pass@www.boot.dev/blog/path",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "with only username",
			rawURL: "https://user@www.boot.dev/blog/path",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "with encoded characters in path",
			rawURL: "https://www.boot.dev/blog/hello%20world",
			want:   "www.boot.dev/blog/hello world",
		},
		{
			name:   "with encoded characters in query",
			rawURL: "https://www.boot.dev/blog/path?q=hello%20world",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "with uppercase host",
			rawURL: "https://WWW.BOOT.DEV/blog/path",
			want:   "WWW.BOOT.DEV/blog/path",
		},
		{
			name:   "empty path with query",
			rawURL: "https://www.boot.dev?key=value",
			want:   "www.boot.dev",
		},
		{
			name:   "ftp scheme",
			rawURL: "ftp://files.example.com/pub/file.txt",
			want:   "files.example.com/pub/file.txt",
		},
		{
			name:   "subdomain",
			rawURL: "https://blog.store.example.com/products",
			want:   "blog.store.example.com/products",
		},
		{
			name:   "ip address host",
			rawURL: "https://192.168.1.1/api/v1",
			want:   "192.168.1.1/api/v1",
		},
		{
			name:   "ipv6 host",
			rawURL: "https://[::1]/api/v1",
			want:   "[::1]/api/v1",
		},
		{
			name:   "path with special chars",
			rawURL: "https://www.boot.dev/path/to/file.name.html",
			want:   "www.boot.dev/path/to/file.name.html",
		},
		{
			name:   "empty string",
			rawURL: "",
			want:   "",
		},
		{
			name:   "double slash scheme",
			rawURL: "https://www.boot.dev//blog/path",
			want:   "www.boot.dev//blog/path",
		},
		{
			name:   "with percent encoded host",
			rawURL: "https://www.boot.dev/blog/path",
			want:   "www.boot.dev/blog/path",
		},
		{
			name:   "complex URL with all components",
			rawURL: "https://user:pass@www.boot.dev:8443/blog/path?foo=bar&baz=qux#section",
			want:   "www.boot.dev:8443/blog/path",
		},
		{
			name:    "mailto scheme",
			rawURL:  "mailto:user@example.com",
			want:    "",
			wantErr: false,
		},
		{
			name:   "file scheme local path",
			rawURL: "file:///etc/passwd",
			want:   "/etc/passwd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := normalizeURL(tt.rawURL)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("normalizeURL() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatalf("normalizeURL() succeeded unexpectedly, got %v with error: %v", got, gotErr)
			}
			if got != tt.want {
				t.Errorf("normalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
