The oauth2 web site is based on gomniauth library. It redirect to google oauth page.

Basic flow:
1. if the request doesn't have valid cookie, redirect the web browser to google oauth URL through the response
2. After login google oauth url, given the client ID, the browser will redirect back to this web site's callback path /callback/ (it's configured in google account)
3. /callback/ path's http handler will get and decode(BASE64) the parameters from callback, complete the auth. It set the cookie back to response and redirect to /
4. the web browser be redirected to / path.  since it has cookie in the request object, the web server parse the user name from cookie and show the page with the name



Followings are google oauth setup:

Google account: GordonJiang1981@hotmail.com, 10 bytes password
 
URL: https://console.developers.google.com/project/gjiangoauth/apiui/credential/oauthclient/263631501515-4l08gqagj93r4iqlk1ciuc1hrpbt69i5.apps.googleusercontent.com
 
Project name:  gjiangoauth
client ID: 263631501515-4l08gqagj93r4iqlk1ciuc1hrpbt69i5.apps.googleusercontent.com
client secrete: oJAYjKrwZlTvLLVSQlWSTKPL
redirect URI http://localhost:8080/callback
