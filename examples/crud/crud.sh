DOCS_DIR=/tmp/docs

create_docs_folder() {
    mkdir -p $DOCS_DIR
}

create_doc() {
    echo $body > $DOCS_DIR/$name
}

read_doc() {
    cat $DOCS_DIR/$name
}

update_doc() {
    echo $body >> $DOCS_DIR/$name
}

delete_doc() {
    rm $DOCS_DIR/$name
}

create_docs_folder
