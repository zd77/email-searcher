## Este repositorio consta de dos programas
1. email-search-app ( VueJS App )
2. indexer ( GoLang Script )

### Instrucciones

1. Levantar contenedor zincsearch con este cmd:
```
docker run -v /full/path/of/data:/data -e ZINC_DATA_PATH="/data" -p 4080:4080 \
    -e ZINC_FIRST_ADMIN_USER=admin -e ZINC_FIRST_ADMIN_PASSWORD=Complexpass#123 \
    --name zincsearch public.ecr.aws/zinclabs/zincsearch:latest
```

* entrar a archivo indexer por terminal
```
cd indexer
```

* ejecutar cmd de construccion de binario
```
    make build
```

* ejecutar indexer:
    - no es necesario incluir las tags < > dentro de ellas basta con incluir un
    path relativo (e.g.) ../../enron_mail_20110402/maildir
```
    ./indexer -path=<direccion de archivo enron_mail_20110402>/maildir
```

* entrar a email-search-app
```
cd email-search-app
```

* generar .env en el archivo raiz y agregar variable
    - VITE_BASE_URL=http://localhost:4080
```
cd email-search-app
```
* npm install
```
npm i
```

* levantar vueJS:
```
npm run dev
```

### IMPORTANT
**una vez arriba el email-search-app tomar en consideracion que en la barra de busqueda basta con agregar algun criterio de busqueda y precionar la tecla enter**

