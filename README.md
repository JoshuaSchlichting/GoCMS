# GoCMS
##### GoCMS is intended to be the boiler plate code you need to build a CMS without having too many dependencies outside of the Go standard library.
___

### Current state of development
GoCMS is in early stages of development and is not even close to being ready for production use. The code in this repository is not stable and will change rapidly over time.

### Dependencies
While the idea is to not have many dependencies, some things are still being imported. At the time of writing this, routing and JWT are handled by [Chi](https://github.com/go-chi/chi) (and they probably always will be).

### Design
As of now, the CMS is relying on free templates from BootStrapMade.com.