# albomx-comics

## Descripción
Aplicación Golang para el desafío de Albomx

La misma realiza una sincronización contra la base de datos de Marvel, a través de su API https://developer.marvel.com/ cada 24 hs, con los heroes Iron Man y Capitan América.
Dicha actualización obtiene los escritores, coloristas y editores de cada comics, junto con los heroes que participaron en los diferenctes comics.

Se ofrecen dos endpoints para ver la respuesta de la sincronización. Uno que obtiene los colaboradores de cada comic y otro que devuelve los heroes que participaron en cada comic.

## Uso remoto
Para hacer uso de los endpoints se puede ingresar a https://test-albo-mx-juan-espinoza.herokuapp.com/swagger/index.html y probar los endpoints con los Ids:
* ironman
* capamerica

Cualquier otro ID retornará un 404 con el mensaje de que no se encutra su ID en la base de datos

## Ejecución local
Si tiene instalado Golang puede hacer uso del archivo Makefile que se encuentra dentro del repositorio, el mismo contiene diferentes funciones para ejecutar la aplicación. Debe ejecutarlas en el sigueinte orden:

* ``dependencies``
* ``run-locally``

Al comando run-locally debe configurarle las siguientes varialbes:
* **DATABASE_URL:** tiene la información de dónde se encuentra la base de datos. **NOTA:** debe ser postgres
* **PRIVATE_KEY:** hace referencia a la KEY privada de su cuenta de Marvel
* **API_KEY:** hace referencia a la KEY pública de su cuenta de Marvel