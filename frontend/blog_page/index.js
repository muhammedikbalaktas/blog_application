var divId=sessionStorage.getItem("blog_id");
sendRequest()
function sendRequest() {
    const baseUrl = 'http://localhost:8080/single_blog';
    const queryParams = { 
        
        blog_id: divId,
    };

    
    const queryString = new URLSearchParams(queryParams).toString();
    const url = `${baseUrl}?${queryString}`;

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
    var path="http://localhost:8080/images/"
    var imageName=data.ImageName;
    var imagePath=path+imageName;
    var id=data.DivId;
    var likeCount=data.LikeCount;
    var likeId="like_"+id;
    var htmlData=`
    <img src="${imagePath}" alt="blog_img">
    <h2>${data.Title}</h2>
    <p>${data.Content}</p>
    
    <div  class="right_bottom">
        <p>Author: ${data.Username}</p>
        <div  class="likes">
            <p id="${likeId}">${likeCount}</p>
            <img src="like_icon.png" alt="like_icon">
            
            
        </div>

    </div>`
    var blogDiv=document.getElementById("single_blog");
    blogDiv.innerHTML=htmlData;
}

