
sendRequest();
function sendRequest(){
    const url = 'http://localhost:8080/get_blogs';



const requestOptions = {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json' 
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
    
    
    parseBlog(data);
    
  })
  .catch(error => {
    console.error('Error:', error);
  });

}
function parseBlog(data){
    for (var i = 0; i < data.length; i++) {
        var obj = data[i];
        addBlog(obj);
    }
    
}
function addBlog(data){
        var path="http://localhost:8080/images/"
        var imageName=data.ImageName;
        var imagePath=path+imageName;
        var croppedContent=data.Content.substring(0,230);
        var id=data.DivId;
        var likeCount=data.LikeCount;
        var likeId="like_"+id;
        var imgId="img_"+id;
        var blogData=`
        <div class="left">
            <img src="${imagePath}" alt="blog_img" class="blog_img">
        </div>
        <div class="right">
            <div class="right_top">
                
                <h2 id="${id}" class="title">${data.Title}</h2>
             
                <p  class="content">${croppedContent}...</p>
            </div>
            <div  class="right_bottom">
                <p>Author: ${data.Username}</p>
                <div  class="likes">
                    <p id="${likeId}">${likeCount}</p>
                    <img src="like_icon.png" alt="like_icon" id="${imgId}">
                    
                    
                </div>

            </div>
          </div>`

        var blogsDiv=document.getElementById("blogs");
        var singleBlog=document.createElement("div")
        
        
        singleBlog.classList.add("blog");
        singleBlog.innerHTML=blogData;
        
        blogsDiv.appendChild(singleBlog);
        


        var title=document.getElementById(id);
        title.addEventListener("click",function(){
          sessionStorage.setItem("blog_id",id);
          window.location.href="http://127.0.0.1:5500/blog_page/"
        })
        var img=document.getElementById(imgId);
        img.addEventListener("click",function(){
          
          
          socket.send(`
          {
              "div_id": "${id}"
          }
          
          `)
        })
}
const userToken=sessionStorage.getItem("token");

const queryParams = {
    token: userToken,
};


const queryString = Object.keys(queryParams).map(key => key + '=' + 
encodeURIComponent(queryParams[key])).join('&');
const socket = new WebSocket(`ws://localhost:8080/ws?${queryString}`);

socket.addEventListener('open', function (event) {
    console.log('WebSocket connected!');
});


socket.addEventListener('message', function (event) {
      try {
        
        var eventData = JSON.parse(event.data);

        
        var divId = eventData.div_id;
        var likeId = "like_" + divId;
        var likeCount = eventData.like_count;
        console.log("value has been set");
        var like=document.getElementById(likeId);
        like.innerText=likeCount;
        
        var hasNotif=eventData.has_notif;
        console.log(hasNotif);
        if(hasNotif){
          var notifIcon=document.getElementById("notification_count");
          notifIcon.style.visibility="visible";
        }

        
    } catch (error) {
        console.error('Error parsing event data:', error);
    }
    // var like=document.getElementById(likeId);
});


socket.addEventListener('close', function (event) {
    console.log('WebSocket connection closed.');
});


socket.addEventListener('error', function (event) {
    console.error('WebSocket error:', event);
});



