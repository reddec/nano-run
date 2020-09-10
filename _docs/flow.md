# High-level overview

Subjects:

* Client - the side which is making HTTP(S) request to the System
* System - instance of nano-run that routing request
* Worker - executable that implements business logic

During restart - all incomplete tasks will queued again.

![image](https://user-images.githubusercontent.com/6597086/92712138-d8b58580-f38b-11ea-8a26-251df5c4ae13.png)

![image](https://user-images.githubusercontent.com/6597086/92578247-3085bb00-f2be-11ea-87de-e2c9d94a21fa.png)

