{{$comment := index .Comments 0 }}
{{$comment.UserComment.StarRating | stars}} {{.AuthorName}}
{{$comment.UserComment.Text}}