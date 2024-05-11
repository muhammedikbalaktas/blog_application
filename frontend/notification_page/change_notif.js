sendRequest();

function sendRequest() {
    var token=sessionStorage.getItem("token");
    const url = new URL('http://localhost:8080/change_notif');
    var token = sessionStorage.getItem("token");

    const params = { 'status': '0' };
    url.search = new URLSearchParams(params).toString();
    
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
        })
        .catch(error => {
            console.error('Error:', error);
        });
}
