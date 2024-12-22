# FastFuzz

Herramienta para fuzzear directorios de un aplicativo web de forma **veloz** gracias a las ventajas que ofrece un lenguaje como **Golang**. Es posible asignarle una cantidad de **Workers** para que en paralelo intenten descubrir directorios y archivos en rutas públicas.

## Instalación

```bash
git clone https://github.com/Yato03/FuzzFast
cd FuzzFast
go build app/fuzzer.go
```

## Uso

### Help

```bash
fuzzer --help
Usage of fuzzer:
  -output string
        Output file
  -t int
        Workers to use (default 100)
  -url string
        URL to fuzz
  -wordlist string
        Wordlist to use
```

Como indica el comando help, existen los siguientes parámetros:

- **url** [obligatorio]: Especifica la url a fuzzear. Deberá de ser de la forma: `http(s)://dominio/`, puediendo ser `http` o `https`.
- **wordlist** [obligatorio]: Fichero donde se determinan los directorios y ficheros a fuzzear en la página web.
- **output** [opcional]: Fichero de salida con los resultados. En él se detallan las rutas encontradas por el fuzzer. Si no se especifica ruta, no generará ningún archivo de output.
- **t** [opcional]: Determina el número de *Workers* que se usarán para paralelizar el trabajo. El número recomendado dependerá de las prestaciones del dispotivo donde se ejecuta. Por defecto serán 100.

### Ejemplo de ejecución

```bash
fuzzer --url http://localhost --wordlist ./wordlist.txt --t 50 --output ./output.txt
```