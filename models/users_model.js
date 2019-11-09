//用户模型

var mongoose=require.resolve('mongoose');
var Schema=mongoose.Schema;
var UserSchema=new Schema({
    username:{type:String,unique:true},
    email:String,
    color:String,
    hashed_password:String
});
mongoose.model('User',UserSchema);