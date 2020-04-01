# psqlisten2websocket

this repository holds code written in go that demonstrates two things.
* database listener that asynchronous notify the client
* websocket server that broadcast messages to all registered clients

*trigger_func.sql*
<pre>
CREATE OR REPLACE FUNCTION notify_event() RETURNS TRIGGER AS $$

    DECLARE
        data json;
        notification json;

    BEGIN

        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
        IF (TG_OP = 'DELETE') THEN
            data = row_to_json(OLD);
        ELSE
            data = row_to_json(NEW);
        END IF;

        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);


        -- Execute pg_notify(channel, notification)
        PERFORM pg_notify('events',notification::text);

        -- Result is ignored since this is an AFTER trigger
        RETURN NULL;
    END;

$$ LANGUAGE plpgsql;
</pre>

*trigger.sql*
<pre>
CREATE TRIGGER products_notify_event
AFTER INSERT OR UPDATE OR DELETE ON portal_eventrequests
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
</pre>

every **INSERT**, **UPDATE** or **DELETE** to postgresql table **portal_eventrequest** will trigger the function notify_event() and these function sends a notification as a JSON string to the listener in *[listener.go](https://github.com/fgau/psqlisten2websocket/blob/master/listener.go)*.

*example JSON string*
<pre>
{
    "table": "portal_eventrequests",
    "action": "INSERT",
    "data": {
        "id": 74036,
        "created_on": "2020-03-30T13:17:41.092166+00:00",
        "url": "https:/veranstaltungen.handelsblatt.com/stadtwerke",
        "termin_id": 26830,
        "user_hash_id": 52,
        "ev_id": 26830,
        "ev_type": 0
    }
}
</pre>

these message was broadcast to all registered clients as you can see here in *[listener.go line 35](https://github.com/fgau/psqlisten2websocket/blob/5616507f8ada57e8f5efe7e865f385e3b3a95353/listener.go#L35)*.

more examples for websocket handling you can find in the official [gorilla websocket repository](https://github.com/gorilla/websocket).

### run and build the server
<pre>
git clone https://github.com/fgau/psqlisten2websocket.git
cd psqlisten2websocket
make run

open a new browser with [http://localhost:8080]()
then add new data to your database table und check the output on your console or browser.
</pre>

feel free to build fancy dashboards, chats and notification systems with go and websockets.
