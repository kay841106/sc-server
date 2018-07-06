
DROP TABLE [Region]
GO

CREATE TABLE [Region] (
[id] int NOT NULL IDENTITY(1,1),
[name] nvarchar(20) COLLATE SQL_Latin1_General_CP1_CI_AS NOT NULL,
[code] nvarchar(10) NULL,
[sort] smallint NOT NULL,
[enable] bit NOT NULL DEFAULT ((0)),
[parent_code] nvarchar(10) NULL,
CONSTRAINT [PK__Region__3213E83F78795B4C] PRIMARY KEY ([id]) ,
CONSTRAINT [uni_code] UNIQUE ([code] ASC)
)
GO


ALTER TABLE [Region] ADD CONSTRAINT [fk_parent_code] FOREIGN KEY ([parent_code]) REFERENCES [Region] ([code]) ON DELETE NO ACTION ON UPDATE NO ACTION
GO
