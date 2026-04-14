interface Props {
  message: string
}

export function ErrorMessage({ message }: Props) {
  return (
    <div role="alert" style={{ color: 'red', padding: '1rem', border: '1px solid red', borderRadius: 4 }}>
      {message}
    </div>
  )
}
