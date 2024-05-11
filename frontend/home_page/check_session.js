const tkn=sessionStorage.getItem("token");


if (tkn==null){
    document.body.innerHTML=`
    <h1>invalid entry. please <a href="http://127.0.0.1:5500/sign_in/index.html">login<a/> again</h1>
    `
}