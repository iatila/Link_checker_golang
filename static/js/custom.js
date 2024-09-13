$(document).ready(function () {
    $("#linkForm").submit(function (event) {
        event.preventDefault();

        let pageLinks = $("#links").val().trim().split("\n").filter(link => link.trim() !== "");
        $("#resultTable tbody").empty();
        $("#resultTable tbody").append('<tr><td colspan="3">Kontrol ediliyor...<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#1c8d2f" d="M12,21L15.6,16.2C14.6,15.45 13.35,15 12,15C10.65,15 9.4,15.45 8.4,16.2L12,21" opacity="0"><animate id="svgSpinnersWifiFade0" fill="freeze" attributeName="opacity" begin="0;svgSpinnersWifiFade1.end+0.2s" dur="0.25s" values="0;1"/><animate id="svgSpinnersWifiFade1" fill="freeze" attributeName="opacity" begin="svgSpinnersWifiFade3.end+0.5s" dur="0.1s" values="1;0"/></path><path fill="#1c8d2f" d="M12,9C9.3,9 6.81,9.89 4.8,11.4L6.6,13.8C8.1,12.67 9.97,12 12,12C14.03,12 15.9,12.67 17.4,13.8L19.2,11.4C17.19,9.89 14.7,9 12,9Z" opacity="0"><animate id="svgSpinnersWifiFade2" fill="freeze" attributeName="opacity" begin="svgSpinnersWifiFade0.end" dur="0.25s" values="0;1"/><animate fill="freeze" attributeName="opacity" begin="svgSpinnersWifiFade3.end+0.5s" dur="0.1s" values="1;0"/></path><path fill="#1c8d2f" d="M12,3C7.95,3 4.21,4.34 1.2,6.6L3,9C5.5,7.12 8.62,6 12,6C15.38,6 18.5,7.12 21,9L22.8,6.6C19.79,4.34 16.05,3 12,3" opacity="0"><animate id="svgSpinnersWifiFade3" fill="freeze" attributeName="opacity" begin="svgSpinnersWifiFade2.end" dur="0.25s" values="0;1"/><animate fill="freeze" attributeName="opacity" begin="svgSpinnersWifiFade3.end+0.5s" dur="0.1s" values="1;0"/></path></svg></td></tr>');

        let completedRequests = 0;
        const totalLinks = pageLinks.length;
        $('#control').attr('disabled',true);
        pageLinks.forEach(function (link) {
            $.ajax({
                url: '/check',
                type: 'POST',
                data: JSON.stringify({ link: link }),
                contentType: "application/json",
                success: function (response) {

                    // Yanıtı kontrol et  ve tablo sınıfını belirle
                    const rowClass = response.status === 0 ? 'table-dark' :
                        response.status >= 200 && response.status < 300 ? 'table-success' :
                            response.status >= 400 && response.status < 500 ? 'table-danger' :
                                response.status >= 300 && response.status < 400 ? 'table-warning' :
                                    response.status >= 500 ? 'table-danger' : '';

                    // Doğru verileri tabloya ekleme
                    $("#resultTable tbody").append(`
                        <tr class="${rowClass}">
                            <td>${response.link || 'Belirtilmemiş'}</td>
                            <td>${response.status !== undefined ? response.status : 'Belirtilmemiş'}</td>
                            <td>${response.description || 'Açıklama mevcut değil'}</td>
                        </tr>
                    `);
                },
                error: function () {
                    // Hata durumunda tabloya ekleme
                    $("#resultTable tbody").append(`
                        <tr class="table-warning">
                            <td>${link}</td>
                            <td>0</td>
                            <td>Bağlantı hatası</td>
                        </tr>
                    `);
                },
                complete: function () {
                    completedRequests++;
                    if (completedRequests === totalLinks) {
                        $('#control').attr('disabled',false);
                        $("#resultTable tbody").find('tr').first().remove();
                    }
                }
            });
        });
    });
});
