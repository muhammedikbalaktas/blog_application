


var userNameInput=document.getElementById("username_input");
var passwordInput=document.getElementById("password_input");

var btnSubmit=document.getElementById("btn_submit")


btnSubmit.addEventListener("click",checkInputs)
userNameInput.addEventListener("input",function(){
    this.style.borderColor="black";
})
passwordInput.addEventListener("input",function(){
    this.style.borderColor="black";
})
function checkInputs(){
    //check if inputs are valid
    if(userNameInput.value.length===0 ||passwordInput.value.length===0){
        
        if(userNameInput.value.length===0 &&passwordInput.value.length===0){
            userNameInput.style.borderColor="red";
            passwordInput.style.borderColor="red";
        }else if(userNameInput.value.length===0){
            userNameInput.style.borderColor="red";
        }else{
            passwordInput.style.borderColor="red";
        }
        return;
        
    }
    sendRequest(userNameInput.value,passwordInput.value)

}


function sendRequest(userNameInput,passwordInput){
    const url = 'http://localhost:8080/get_user';
    const data = {
        username: userNameInput,
        password: passwordInput
    };
    

const requestOptions = {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json' 
  },
  body: JSON.stringify(data) 
};


fetch(url, requestOptions)
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json(); 
  })
  .then(data => {
    
    const token=data["success"]
    sessionStorage.setItem("token",token);
    window.location.href="http://127.0.0.1:5500/home_page/index.html"
    
  })
  .catch(error => {
    console.error('Error:', error);
  });

}