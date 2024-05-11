const token=sessionStorage.getItem("token")

window.addEventListener('hashchange', function(event) {
  
  console.log('User returned to this page');
  
  
  sendRequest();
});
if (token===null){
    console.log("invalid token");
    
}else{
    sendRequest();
}
function sendRequest(){
    const url = 'http://localhost:8080/get_notif';



const requestOptions = {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'authorization': token
  }
};


fetch(url, requestOptions)
  .then(response => {
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }
    return response.json(); 
  })
  .then(data => {
    
    console.log(data);
    changeUI(data);
    
  })
  .catch(error => {
    console.error('Error:', error);
  });

}

function changeUI(data){
    if(data.success){
        var notificationUI =document.getElementById("notification_count");
        notificationUI.style.visibility="visible";

    }else{
        var notificationUI =document.getElementById("notification_count");
        notificationUI.style.visibility="hidden";
    }
}