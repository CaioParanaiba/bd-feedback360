DROP TABLE IF EXISTS feedback;
DROP TABLE IF EXISTS usuario;


CREATE TABLE usuario (
    id_usuario INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    tipo_usuario VARCHAR(30) NOT NULL CONSTRAINT ck_tipo_usuario CHECK (tipo_usuario IN (
                                                                    'porteiro', 'limpeza',
                                                                    'supervisor', 'morador'))
);

CREATE TABLE feedback (
    id_feedback INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    descricao TEXT NOT NULL,
    data_hora TIMESTAMP NOT NULL DEFAULT now(),
    status VARCHAR(30) NOT NULL CONSTRAINT ck_status_feedback CHECK (status IN (
                                                                    'registrado','revisado','arquivado')),
    id_emissor INT NOT NULL CONSTRAINT fk_feedback_emissor REFERENCES usuario(id_usuario),
    id_receptor INT NOT NULL CONSTRAINT fk_feedback_receptor REFERENCES usuario(id_usuario),

    CONSTRAINT ck_feedback_emissor_receptor CHECK (id_emissor <> id_receptor)
);