# WishTrack

A ideia é ser um app onde o usuário possa cadastrar livros/filmes/jogos e conteúdos em geral que ele deseja consumir ou já consumiu. Em relação aos que ele ainda deseja consumir ele poderá indicar o quanto quer, e em relação aos que já consumiu ele poderá avaliar. Em resumo é praticamente um "Ver depois" + "favoritos" 🤷‍♂️


## Pensando as features (O que diabos eu quero 🤔)

* CRUD de conteúdo (se pá as categorias e generos podem ser cadastrados juntos aqui).
* Ter um feed onde serão exibidos alguns conteúdos sujeridos (pensar em um algoritmo legal aqui para sugerir as paradas com algum sentido).
* Marcar um conteúdo como visto e dar uma avaliação.


### Entidade Content (Conteúdo)
```json
{
    "name": "Senhor dos aneis: O retorno do rei",
    "category": "book",
    "genres": [
        "fantasy",
        "medieval",
        "adventure"
    ],
    "summary": "The Return of the King, the third and final volume in The Lord of the Rings",
    "wish_level": 5, // 1 - 5 o quanto eu quero ver isso

    // isso aqui vai identificar o conteudo como visto ✅
    "rate": 5, // 1 - 5 o quanto eu curti OBRIGATÓRIO - indica como visto
    "comment": "Um clássico" // meu comentário
}
```