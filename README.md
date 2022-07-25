# realtime database

Something like Firestore where you could connect to a websocket from clientside
and then any updates from REST would be streamed in form of events. The schema
is Users to start with but could probably be applied to any sort of JSON
object-based data.

Here ArangoDB is used.

## Setup

Required env vars (in the `.env` file):

```env
ARANGO_ROOT_PASSWORD string
DB_PROTOCOL          string
DB_HOST              string
DB_PORT              string
DB_NAME              string
DB_COLLECTION        string
```

To run:

```bash
docker-compose up --build
``s`
