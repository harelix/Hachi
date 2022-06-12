
### Server configuration
dynamic handlers and routing 

```json
{
  "repository" : [
    {
      "pattern" : "/events_history/",
      "query" : "select * from snitch_messages_flow where elementId={{.key}}",
      "type" : "query",
      "sourceAlias" : "postgres1"
    },
    {
      "pattern" : "/history_old/{key}",
      "query" : "select array_to_json(array_agg(row_to_json(t)))  from (SELECT  time, origin,(outgoing_raw_message->'Metadata') as metadata,outgoing_raw_message->'Metadata'->'Executor' as Ex, outgoing_raw_message->'Metadata'->'Event' as event,outgoing_raw_message->'Metadata'->'Executor' as Executor, SPLIT_PART(event,'.', 1) as base, SPLIT_PART(event,'.', 2) as predicate, SPLIT_PART(event,'.', 3) as object, outgoing_raw_message->'ElementId' as ElementId,outgoing_raw_message->'PropertyIdValue' as Property FROM snitch_messages_flow  WHERE outgoing_raw_message->>'ElementId' = '{{.key}}' ORDER BY time DESC)t ",
      "type" : "query",
      "sourceAlias" : "postgres1"
    },
    {
      "pattern" : "/history/{{element_id}}/{{order_by}}",
      "query" : "select array_to_json(array_agg(row_to_json(t)))  from (SELECT * FROM public.history_event where element_id='{{.element_id}}' order by {{.order_by}} desc)t ",


      "type" : "query",
      "sourceAlias" : "postgres1"
    },
    {
      "pattern" : "/history/{key}/{id:[0-9]+}",
      "query" : "",
      "type" : "sp",
      "sourceAlias" : "sql1"
    }
  ]
}
```