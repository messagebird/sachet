# websms.com (Link Mobility)

This sachet provider uses the *websms* SMS gateway service (websms.com). This service is popular in Austria and other European countries.

See https://www.websms.com/ and https://developer.websms.com/web-api/


## Configuration

This provieder needs an API key to access the REST API. To create one login to your websms account got to "API administration" and "API access data". There you can create a new API token for use in this provider.

Note: use a dedicate API token for access by sachet. **DONT use tokens created for productive systems and services**.

You can specifiy a list of target phone numbers which will reveive messages. You have to specificy those numbers in E.164/MSISDN format, e.g. "49123456789" (otherwise often written as "+49123456789" or "0049123456789"). See https://en.wikipedia.org/wiki/MSISDN