mongosh -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" admin <<EOF
db = db.getSiblingDB("$MONGO_INITDB_NAME");
db.createUser({
'user': "$MONGO_NEWUSER_NAME",
'pwd': "$MONGO_NEWUSER_PASSWORD",
'roles': [
      {'role': 'dbOwner', 'db': "$MONGO_INITDB_NAME"}
   ]
});
db.createCollection("$MONGO_INITDB_COL_USER", {
   validator: {
      \$jsonSchema: {
         bsonType: "object",
         required: [ "email", "passHash", "confirmed" ],
         properties: {
            email: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            passHash: {
               bsonType: "binData",
               description: "must be a byte array and is required"
            },
            confirmed: {
               bsonType: "bool",
               description: "must be a bool and is required"
            },
            admin: {
               bsonType: "bool",
               description: "must be a bool if the field exists"
            },
         }
      }
   }
});
db.$MONGO_INITDB_COL_USER.createIndex({ email: 1 },{ unique: true });
db.createCollection("$MONGO_INITDB_COL_APP", {
   validator: {
      \$jsonSchema: {
         bsonType: "object",
         required: [ "name", "secret" ],
         properties: {
            name: {
               bsonType: "string",
               description: "must be a string and is required"
            },
            secret: {
               bsonType: "string",
               description: "must be a string and is required"
            },
         }
      }
   }
});
EOF