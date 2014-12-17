snapchat-spam
===
---

Fucking advanced app written in go. Call it like so:
> $ go run .\app.go [snapchat-access-token] [victim-username] [your-username] [number-of-snaps-to-send]

For example:
> $ go run .\app.go 68925bb39093a89279762446df611f64 alexerax notshad 1000

It will pull in the jpg image specified in app.go as `BaePath` and sends it x number of times to the recepient 
you specified. Currently the specified jpg is broken and snapchat refuses to load it in the app, causing a 
notification that you can't clear without clearing the entire conversation in settings. So backup those nudes 
before you use it on your loved one.

### Getting your Snapchat AuthToken
Close the app - force close it - attach a network insector to your phone (I recommend Charlex Proxy or Fiddler). 
Open the app on your phone and look at the request made to `https://feelinsonice-hrd.appspot.com/loq/all_update` 
in the json response, under `updates` there is a field labeled `auth_token`. Copy that and you're all gucci.


