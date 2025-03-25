BEGIN TRANSACTION;

INSERT INTO publisher (name)
VALUES
('pagina12'),
('perfil'),
('lapolitica');

INSERT INTO message
(status_id, publisher_id, content)
VALUES
(0, 1, 'some post from Pagina 12 that already Done'),
(2, 3, 'some post from La Politica that processing right now'),
(1, 1, 'some post from Pagina 12'),
(1, 1, 'some other post from Pagina 12'),
(1, 2, 'some post from Perfil'),
(1, 3, 'some post from La Politica');

COMMIT;