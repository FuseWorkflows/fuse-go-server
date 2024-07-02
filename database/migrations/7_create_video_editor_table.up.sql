CREATE TABLE video_editor (
  video_id UUID NOT NULL REFERENCES videos(id),
  editor_id UUID NOT NULL REFERENCES editors(id),
  PRIMARY KEY (video_id, editor_id)
);