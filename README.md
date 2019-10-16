# Rabbit Producer
Simple CLI utility for publishing messages to RabbitMQ

## Usage

Grab the appropriate binary (Linux or Mac) from the [dist](dist) folder.

### Basic Usage

```
rabbit-producer -remoteUrl=[user:pass@]myrabbit.com/virtualhost -exchange=profile -routingKey=new-image -message=HI
```

You can also pipe input from stdin:

```
cat message.json | rabbit-producer -remoteUrl=[user:pass@]myrabbit.com/virtualhost -exchange=profile -routingKey=new-image
```

#### Content Types
If rabbit-producer can tell your message input is JSON, it will automatically set a content-type header of `application/json`. If you aren't ok with that, use the `-contentType=` flag to override it.

### Shameless Self Promotion
Want to see the messages on the other side? Check my other project, [rabbit-tail](https://github.com/ricanontherun/rabbit-tail)

# TODO
1. tests