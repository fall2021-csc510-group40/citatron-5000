db.createUser(
    {
        user: "core",
        pwd: "password",
        roles: [
            {
                role: "readWrite",
                db: "citeman"
            }
        ]
    }
);

db.createCollection("works", {
    validator: { $jsonSchema: {
        bsonType: ["object"],
        required: [
            "title",
            "authors",
            "hash",
        ],
        properties: {
            title: {
                bsonType: ["string"],
                minLength: 1,
                description: "must be a non-empty string and is required",
            },
            authors: {
                bsonType: ["array"],
                minItems: 1,
                description: "must be a non-empty array of strings and is required",
                items: {
                    bsonType: ["string"],
                    minLength: 1,
                    description: "must be a non-empty string"
                }
            },
            hash: {
                bsonType: ["string"],
                minLength: 1,
                description: "must be a unique non-empty string"
            }
        }
    }}
});

works = db.getCollection("works");
works.createIndex( { title: "text" } );
works.createIndex( { hash: 1 }, { unique: true } );

