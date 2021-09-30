# REST API

If an `error` field is present, then the rest of the data will not be defined.

## GET /search

Performs a search for a work via its properties. Currently, only title is supported.

### Body (JSON)

```
{
    "title": <TITLE>
}
```

### Response (JSON)

```
{
    "results": [
        "id": <STRING>,
        "hash": <STRING>,

        "type": <STRING>,

        "doi": <STRING>,
        "arxiv": <STRING>,
        "isbn": <STRING>,

        "title": <STRING>
        "authors": [
            <AUTHOR1>,
            <AUTHOR2>,
            ...
        ],

        "version": <STRING>,
        "venue": <STRING>,
        "page": <STRING>,

        "year": <NUMBER>,
        "month": <NUMBER>,
        "day": <NUMBER>,

        "keywords": [
            <KEYWORD1>,
            <KEYWORD2>,
            ...
        ]
    ],
    error: <STRING>
}
```

## GET /format

### Body (JSON)

```
{
    "id": <STRING>
}
```

or

```
{
    "type": <STRING>,

    "doi": <STRING>,
    "arxiv": <STRING>,
    "isbn": <STRING>,

    "title": <STRING>
    "authors": [
        <AUTHOR1>,
        <AUTHOR2>,
        ...
    ],

    "version": <STRING>,
    "venue": <STRING>,
    "page": <STRING>,

    "year": <NUMBER>,
    "month": <NUMBER>,
    "day": <NUMBER>,

    "keywords": [
        <KEYWORD1>,
        <KEYWORD2>,
        ...
    ]
}
```

### Response (JSON)

```
{
    result: <STRING>,
    error: <STRING>
}
```