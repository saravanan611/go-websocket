var Gurl = "";
var Gconnection = "";
var Gaction = "";
const userId = Math.floor(Math.random() * 10000).toString();


function extractHost(url) {
    // Create a URL object
    const urlObject = new URL(url);

    // Extract the host (domain + port)
    const host = urlObject.host;

    return host;
}


function showAlert(message, type = "error") {
    const alertDiv = document.createElement("div");
    alertDiv.textContent = message;
    alertDiv.style.position = "fixed";
    alertDiv.style.top = "20px";
    alertDiv.style.left = "50%";
    alertDiv.style.transform = "translateX(-50%)";
    alertDiv.style.padding = "10px 20px";
    alertDiv.style.borderRadius = "5px";
    alertDiv.style.color = "white";
    alertDiv.style.backgroundColor = type === "error" ? "#dc3545" : "#28a745";
    alertDiv.style.zIndex = "1000";

    document.body.appendChild(alertDiv);

    setTimeout(() => {
        alertDiv.remove();
    }, 3000);
}

async function Senddata(url = '', data = {}) {
    // Default options are marked with *
    console.log("data", data,"Gconnection:",Gconnection,"Gaction:",Gaction);

    const response = await fetch(url, {
        method: 'POST', // HTTP method
        headers: {
            'Content-Type': 'application/json', // Set the content type to JSON
            'action': Gaction,
            'conn': Gconnection,
        },
        body: JSON.stringify(data), // Convert data to JSON string
    });
    
    console.log("response",response);
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    }
    return response.json();
}