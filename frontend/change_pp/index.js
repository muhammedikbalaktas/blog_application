const token=sessionStorage.getItem("token");
var button=document.getElementById("submit");
button.addEventListener("click",sendRequest);
console.log(token);
function sendRequest(){
    
    var url = "http://localhost:8080/upload_image";
    var imageFileInput = document.getElementById("image");

    
    

    
    var formData = new FormData();

    
    formData.append("image", imageFileInput.files[0]);

    

    fetch(url, {
        method: "POST",
        body: formData,
        headers: {
            'authorization': token
        }
    })
    .then(response => {
        console.log(response);
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        
        return response.json(); 
    })
    .then(data => {
       
        window.location.href="http://127.0.0.1:5500/home_page/index.html"
    })
    .catch(error => {
        
        console.error('There was a problem with the fetch operation:', error);
    });

}