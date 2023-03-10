package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipReadCloser struct {
	io.ReadCloser
	Reader io.Reader
}

func (r gzipReadCloser) Read(b []byte) (int, error) {
	return r.Reader.Read(b)
}

//GzipDecompressHandle middleware обработчик подменяет reader на gzip.reader
//если клиент посылает сжатый запрос
func GzipDecompressHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что тело запроса псодержит gzip-сжатие
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Reader поверх текущего Body
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		r.Body = gzipReadCloser{Reader: gz, ReadCloser: r.Body}
		//передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(w, r)
	})
}
