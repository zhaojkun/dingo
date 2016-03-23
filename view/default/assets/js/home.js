$(document).ready(function () {
    //fixHeader();
    // topButton();
    renderMarkdown();
    initComment();
});

function fixHeader() {
    var $nav = $('#main-nav');
    var top = $nav.offset().top;
    console.log(top);
    $(window).scroll(function () {
        if (top < $(this).scrollTop()) {
            $nav.addClass("fixed").removeClass("text-center");
        } else {
            $nav.removeClass("fixed").addClass("text-center");
        }
    });
}

// function topButton() {
//     var top = $('#main-nav').offset().top;
//     var $top = $('#go-top');
//     $(window).scroll(function () {
//         if (top < $(this).scrollTop()) {
//             $top.removeClass("hide");
//         } else {
//             $top.addClass('hide');
//         }
//     });
//     $top.on("click", function () {
//         $('body,html').animate({scrollTop: 0}, 500);
//         return false;
//     })
// }

function renderMarkdown() {
    var $md = $('.markdown');
    $md.each(function (i, item) {
        $(item).html(marked($(item).html().replace(/&gt;/g, '>')));
    });
    var code = $md.find('pre code');
    if (code.length) {
        $("<link>").attr({ rel: "stylesheet", type: "text/css", href: "/static/css/highlight.css"}).appendTo("head");
        $.getScript("/static/lib/highlight.min.js", function () {
            code.each(function (i, item) {
                hljs.highlightBlock(item)
            });
        });
    }
}

function initComment() {
    var $list = $('#comment-list');
    if (!$list.length) {
        return;
    }
    if (localStorage.getItem("comment-author")) {
        $('#comment-author').val(localStorage.getItem("comment-author"));
        $('#comment-email').val(localStorage.getItem("comment-email"));
        $('#comment-url').val(localStorage.getItem("comment-url"));
        $('#comment-avatar').attr("src", localStorage.getItem("comment-avatar"));
        $('.c-avatar').removeClass("null");
    }
    $('#comment-content').on("focus", function () {
        if ($('.c-avatar').hasClass("null")) {
            $('.c-avatar-field').remove();
            $('.c-info-fields').removeClass("hide");
        }
    });
    // $('.not-me').on("click", function () {
    //     $('.c-avatar-field').remove();
    //     $('.c-info-fields').removeClass("hide");
    //     return false;
    // });
    $('#comment-show').on("click", function () {
        $('#comment-show').hide();
        $('#comment-form').removeClass("hide");
    });
    $('#comment-cancel').on("click", function () {
        $('#comment-form').addClass("hide");
        $('#comment-show').show();
    });
    $('#comment-form').ajaxForm(function (json) {
      console.log(json);
        if (json.res) {
            localStorage.setItem("comment-author", $('#comment-author').val());
            localStorage.setItem("comment-email", $('#comment-email').val());
            localStorage.setItem("comment-url", $('#comment-url').val());
            localStorage.setItem("comment-avatar", json.comment.avatar);
            var tpl = $($('#comment-tpl').html());
            tpl.find(".comment-avatar").attr("src", json.comment.avatar).attr("alt", json.comment.avatar);
            tpl.find(".comment-name").attr("href", json.comment.website).text(json.comment.author);
            tpl.find(".comment-reply").attr("rel", json.comment.id);
            tpl.find(".comment-content").html("<p>" + json.comment.content + "</p>");
            var date = new Date(json.comment.create_time);
            tpl.find(".comment-message").html("Your comment is awaiting moderation.");
            if (json.comment.parent_md) {
                tpl.find(".c-p-md").html(marked(json.comment.parent_md));
            } else {
                tpl.find(".c-p-md").remove();
            }
            tpl.attr("id", "comment-" + json.comment.id);
            if (json.comment.status == "approved") {
                tpl.find(".comment-check").remove();
            }
            $list.append(tpl);
            $('.cancel-reply').trigger("click");
            $('#comment-content').val("");
        } else {
            alert("Can not submit comment!");
        }
    });
    $list.on("click", ".comment-reply", function () {
        var id = $(this).attr("rel");
        var pc = $('#comment-' + id);
        var md = "> @" + pc.find(".comment-name").text() + "\n\n";
        md += "> " + pc.find(".comment-content").html() + "\n";
        $('#comment-reply').html(marked(md));
        $('#comment-show').hide();
        $('#comment-form').removeClass("hide");
        $('.cancel-reply').show();
        var top = $('#comment-form').offset().top;
        $('body,html').animate({scrollTop: top}, 500);
        return false;
    });
    $('.cancel-reply').on("click", function () {
        $('#comment-reply').empty();
        $('#comment-parent').val(0);
        $(this).hide();
        return false;
    });
}
