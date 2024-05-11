var path="http://localhost:8080/images/"


sendRequest();
function sendRequest(){
    const url = 'http://localhost:8080/get_pp';
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
    
    
    console.log(data);
    addImage(data);
    
  })
  .catch(error => {
    console.error('Error:', error);
  });

}

function addImage(data){
    var path="http://localhost:8080/ppimages/"
    var imageName=data.success;
    var imagePath=path+imageName;
   
    var ppImg=document.getElementById("pp");
    ppImg.src=imagePath;
    
}