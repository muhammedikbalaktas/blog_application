var btnSignIn=document.getElementById("sign_in_btn");


btnSignIn.addEventListener("click",function(){
    window.location.href='http://127.0.0.1:5500/sign_in/index.html'
})
document.getElementById("btn_go").addEventListener("click", function(event) {
    var email = document.getElementById("email").value;
    
    var validRegex = /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/;
    if (email.match(validRegex)) {
      alert("valid email address");
      event.preventDefault();
    }else{
        alert("invalid email")
    }
  });

