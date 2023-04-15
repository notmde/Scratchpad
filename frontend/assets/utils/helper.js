function dropDown() {
  var x = document.getElementsByClassName("dialog-box")[0].classList.toggle("show-box");
}

function myFunction() {
  var x = document.getElementById("password");
  if (x.type === "password") {
    x.type = "text";
  } else {
    x.type = "password";
  }
}