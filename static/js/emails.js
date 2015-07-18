$(document).ready(function() {
    $('#emailToOperate').focus();
    $('#submit_btn').on('click', (function() {
        option = $('#emailOption').val()
        email = $('#emailToOperate').val();
        if (option == "add") {
            $.post("/addemail?email="+email);
        }
        else if (option == "remove") {
            $.post("/removeemail?email=" + email);
        }
        window.location = 'emails.html'
    }));
})