function deleteUrl() {
    if (confirm('Delete this record?')) {
        let id  = event.target.parentNode.parentNode.id
        sendDelete('/deleteurl/', id)
    } 
}

function sendDelete(path, id) {
    const form = document.createElement('form');
    form.method = 'post';
    form.action = path;
  
    const hiddenField = document.createElement('input');
    hiddenField.type = 'hidden';
    hiddenField.name = "ID";
    hiddenField.value = id;

    form.appendChild(hiddenField);
    document.body.appendChild(form);
    form.submit();
  }
  