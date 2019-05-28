$(document).ready(function () {

    $(".deleteBlog").on("click", function () {
        if (confirm("Are you sure want to delete this blog?")) {
            var id = $(this).attr('data-value');
            $.ajax({
                url: "/blog/delete/" + id,
                type: "GET",
                dataType: 'JSON',
                success: function (data) {
                    alert(data.msg);
                    window.location.reload();
                }
            });
        }
    });

});