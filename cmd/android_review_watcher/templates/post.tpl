{{$comment := index .Comments 0 }}
{{$comment.UserComment.StarRating | stars}}
{{$comment.UserComment.Text}}
