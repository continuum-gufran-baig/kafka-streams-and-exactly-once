# What is the Idempotent Producer feature?

When a producer sends messages to a topic, things can go wrong, such as short connection failures. When this happens, any messages that are pending acknowledgements can either be resent or discarded. The messages may have been successfully written to the topic, or not, there is no way to know. If we resend then we may duplicate the message, but if we don’t resend then the message may essentially be lost.


Also, resending messages can cause message ordering to go wrong. So we can end-up with messages delivered twice and out-of-order. The idempotent producer feature addresses these issues ensuring that messages always get delivered, in the right order and without duplicates.


### Limitation 1: Acks=All
For one, you can only use acks=all.

### Limitation 2: max.in.flight.requests.per.connection <= 5
You must either leave max.in.flight.requests.per.connection (alias max.in.flight) left unset so that it be automatically set for you or manually set it to 5 or less

# How it Works

When enable.idempotence is set to true, no manual retries are required, in fact performing retries in your application code will still cause duplicates. Leave the retries to your client library, it is totally transparent to you as a developer.


So retries are taken of, but how can the broker identify duplicate messages and discard them?


Each producer gets assigned a Producer Id (PID) and it includes its PID every time it sends messages to a broker. Additionally, each message gets a monotonically increasing sequence number. A separate sequence is maintained for each topic partition that a producer sends messages to. On the broker side, on a per partition basis, it keeps track of the largest PID-Sequence Number combination is has successfully written. When a lower sequence number is received, it is discarded.



