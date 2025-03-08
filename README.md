# WishTrack (Em desenvolvimento...)

A ideia é ser um app onde o usuário possa cadastrar livros/filmes/jogos e conteúdos em geral que ele deseja consumir ou já consumiu. Em relação aos que ele ainda deseja consumir ele poderá indicar o quanto quer, e em relação aos que já consumiu ele poderá avaliar. Em resumo é praticamente um "Ver depois" + "favoritos" 🤷‍♂️


## Executando este projeto

### **Requisitos**:
 * [Docker/Docker compose](https://www.docker.com/)
 * [Air](https://github.com/air-verse/air)
 * [Migrate](https://github.com/golang-migrate/migrate)
 * [Swag](https://github.com/swaggo/swag)
 * [GNU Make](https://www.gnu.org/software/make/manual/make.htm) *(opcional, mas conveniente)*

#### **Executando projeto local (com hot-reload)**
```shell
    make dev
```

#### **Rodando testes**
```shell
    make tests
```

#### **Gerando executável**
```shell
    make build
```