sendRequest();
function sendRequest(){
    const url = 'http://localhost:8080/get_notifications';
    var token=sessionStorage.getItem("token");


const requestOptions = {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
    'authorization':token
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
    
    
    parseNotification(data);
    
  })
  .catch(error => {
    console.error('Error:', error);
  });

}
function parseNotification(data){
    
    var baseDiv=document.getElementById("home");
    baseDiv.innerHTML='';
    var returnDiv=document.createElement("div");
    returnDiv.classList.add("return");
    returnDiv.innerHTML=`
    <a href="http://127.0.0.1:5500/home_page/index.html">
    <img src="back.png" alt="Description of the image">
    </a>

    `
    baseDiv.appendChild(returnDiv);
    for (let index = 0; index < data.length; index++) {
        const element = data[index];
        var div=document.createElement("div");
        div.classList.add("notification");
        var pTag=document.createElement("p");
        pTag.innerText=element;
        div.appendChild(pTag);
        baseDiv.appendChild(div);

    }
}