-- script MySql, ejecutarlo en el workbenck

CREATE DATABASE `goDB`;

-- Empezamos a usarla en MySQL
USE `goDB`;

-- Tabla de estudiantes
CREATE TABLE IF NOT EXISTS `goDB`.`Students` (
  `id` INT NOT NULL AUTO_INCREMENT, -- ID de la tabla 
  `nombre` VARCHAR(50) NULL, -- Nombre del lenguaje de programación
  `edad` VARCHAR(5) NULL, -- Autor del lenguaje
  PRIMARY KEY (`id`));

-- Insertamos dos filas

-- Primera fila
INSERT INTO `goDB`.`Students`
(`nombre`,
`edad`)
VALUES (
'Jose Issac',
'22');

-- Segunda fila
INSERT INTO `goDB`.`Students`
(`nombre`,
`edad`)
VALUES (
'Omar',
'24');