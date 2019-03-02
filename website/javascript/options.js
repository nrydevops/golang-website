function validatePasswordChange() {

    removeMessages();

    var oldPass = document.forms["change_password"]["old_password"].value;
    var newPass = document.forms["change_password"]["new_password"].value;
    var verifyPass = document.forms["change_password"]["verify_password"].value;

    if (oldPass.length < 3) {

        errorMessage('Invalid old password!');

        return false;
    }
    else if(newPass.length < 3) {

        errorMessage('Invalid new password!');

        return false;
    }
    else if(newPass !== verifyPass) {

        errorMessage('Verification failed!');

        return false;
    }

    return true;
}

function errorMessage(message) {

    var passwordForm = document.getElementById("change_password");
    var paragraphElement = document.createElement('p');
    var italicElement = document.createElement('i');

    italicElement.appendChild(document.createTextNode(message))
    paragraphElement.appendChild(italicElement);
    paragraphElement.setAttribute("class", "wrong_input");

    passwordForm.appendChild(paragraphElement);
}

function removeMessages() {

    var wrongInputElements = document.getElementsByClassName("wrong_input");
    var successElements = document.getElementsByClassName("success");

    while(wrongInputElements.length > 0) {

        wrongInputElements[0].parentNode.removeChild(wrongInputElements[0]);
    }

    while(successElements.length > 0) {

        successElements[0].parentNode.removeChild(successElements[0]);
    }

}