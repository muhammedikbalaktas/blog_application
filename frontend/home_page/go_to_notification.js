var notificationDiv=document.getElementById("notification");


notificationDiv.addEventListener("click", function(){

    window.location.href="http://127.0.0.1:5500/notification_page/index.html";
    
})
var logOut=document.getElementById("log_out");

logOut.addEventListener("click", function(){

    sessionStorage.clear();
    window.location.href="http://127.0.0.1:5500/sign_in/index.html"


})