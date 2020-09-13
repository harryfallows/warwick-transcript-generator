const inputFile0 = document.querySelector(
    "#inputFile0 input[type=file]"
);
const inputFile1 = document.querySelector(
    "#inputFile1 input[type=file]"
);
const inputFile2 = document.querySelector(
    "#inputFile2 input[type=file]"
);
const inputFile3 = document.querySelector(
    "#inputFile3 input[type=file]"
);
inputFile0.onchange = () => {
    if (inputFile0.files.length > 0) {
        const fileName = document.querySelector(
            "#inputFile0 .file-name"
        );
        fileName.textContent = inputFile0.files[0].name;
    }
};
inputFile1.onchange = () => {
    if (inputFile1.files.length > 0) {
        const fileName = document.querySelector(
            "#inputFile1 .file-name"
        );
        fileName.textContent = inputFile1.files[0].name;
    }
};
inputFile2.onchange = () => {
    if (inputFile2.files.length > 0) {
        const fileName = document.querySelector(
            "#inputFile2 .file-name"
        );
        fileName.textContent = inputFile2.files[0].name;
    }
};
inputFile3.onchange = () => {
    if (inputFile3.files.length > 0) {
        const fileName = document.querySelector(
            "#inputFile3 .file-name"
        );
        fileName.textContent = inputFile3.files[0].name;
    }
};

submitForm.onclick = function() {
    var inputs, index;
    var dict = {};
    dict["files"] = {};

    inputs = document.getElementsByTagName("input");

    for (index = 0; index < inputs.length; ++index) {
        currentInput = inputs[index];
        console.log(currentInput.type);
        if (currentInput.type == "file") {
            var file = currentInput.files[0];
            if (file) {
                var reader = new FileReader();
                reader.readAsText(file);
                reader.onload = function(fileContents) {
                    dict["files"][currentInput.name] = fileContents.target.result;
                }
            }
        }
        else if (currentInput.type == "radio") {
            if (currentInput.checked) {
                if (currentInput.name == "logo") {
                    dict["logo"] = currentInput.value;
                }
            }
        }
        else if (currentInput.type == "text") {
            dict[currentInput.id] = currentInput.value;
        }
        else if (currentInput.type == "checkbox") {
            dict[currentInput.name] = currentInput.checked;
        }

    }
    console.log(dict)
    var request = new XMLHttpRequest();
    request.open("POST", "http://localhost:8080", true);
    request.setRequestHeader('Content-Type', 'application/json');
    console.log(JSON.stringify(dict))
    request.send(JSON.stringify(dict));
}

